package main

import (
    "os"
    "flag"
    "fmt"
    "strings"
    "runtime"
    "github.com/go-gl/glfw/v3.2/glfw"
    "github.com/go-gl/gl/v2.1/gl"
    "github.com/fogleman/gg"
    "image"
    "image/draw"
    "bufio"
    "cnc-tools-go/cnclib"
    "github.com/Succo/wkttoorb"
    "github.com/paulmach/orb"
)

type Mesh struct {
    program_id uint32
    vertex_id uint32
    uv_id uint32
    tex_id uint32
    nr_vertices uint32
}

type Context struct {
    Width int
    Height int
    Maxx int
    Maxy int
    Scale float64
}

func main() {
    runtime.LockOSThread()

    var ctx Context
    ctx.Width = 640
    ctx.Height = 480
    ctx.Maxx = 200
    ctx.Maxy = 290
    ctx.Scale = 1.0

    flag.IntVar(&ctx.Width, "width", 640, "Window width")
    flag.IntVar(&ctx.Height, "height", 480, "Window height")
    flag.IntVar(&ctx.Maxx, "maxx", 200, "Maximum X travel on machine")
    flag.IntVar(&ctx.Maxy, "maxy", 290, "Maximum Y travel on machine")
    flag.Parse()

    var things []orb.Geometry
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        geo, err := wkttoorb.Scan(scanner.Text())
        if err != nil {
            panic(err)
        }
        things = append(things, geo)
    }
    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
    }
    lss := cnclib.GeometryToLineStrings(things)

    if err := glfw.Init(); err != nil {
        panic(err)
    }
    defer glfw.Terminate()

    window, err := glfw.CreateWindow(ctx.Width, ctx.Height, "test", nil, nil)
    if err != nil {
        panic(err)
    }

    window.MakeContextCurrent()

    if err := gl.Init(); err != nil {
        panic(err)
    }

    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Fprintln(os.Stderr, "OpenGL version", version)

    initGL()
    gl.Viewport(0, 0, int32(ctx.Width), int32(ctx.Height))
    mesh := generateBuffers()
    drawLineStrings(lss, &ctx)

    window.SetSizeCallback(func(w *glfw.Window, nw int, nh int) {
        ctx.Width = nw
        ctx.Height = nh
        gl.Viewport(0, 0, int32(ctx.Width), int32(ctx.Height))
        drawLineStrings(lss, &ctx)
    })

    window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
        if yoff > 0 {
            ctx.Scale += 0.1
            drawLineStrings(lss, &ctx)
        }
        if yoff < 0 {
            ctx.Scale -= 0.1
            drawLineStrings(lss, &ctx)
        }
    })

    for (!window.ShouldClose()) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
        gl.DrawArrays(gl.TRIANGLES, 0, int32(mesh.nr_vertices))
        glfw.PollEvents()
        window.SwapBuffers()
    }
}

func drawLineStrings(lss []orb.LineString, ctx *Context) {
    dc := gg.NewContext(ctx.Width, ctx.Height)
    dc.ScaleAbout(ctx.Scale, ctx.Scale, 0, 0)
    var colors = [][]float64{
        {1.0, 0.0, 0.0},
        {0.0, 1.0, 0.0},
        {0.0, 0.0, 1.0},
    }
    var c uint32 = 0
    dc.SetRGB(colors[c][0], colors[c][1], colors[c][2])
    for i := 0; i < len(lss); i++ {
        dc.MoveTo(lss[i][0][0], lss[i][0][1])
        for j := 1; j < len(lss[i]); j++ {
            dc.LineTo(lss[i][j][0], lss[i][j][1])
        }
        dc.Stroke()
        c++
        if c > 2 {
            c = 0
        }
        dc.SetRGB(colors[c][0], colors[c][1], colors[c][2])
    }
    dc.SetRGB(0.6, 0.6, 0.6)
    dc.SetDash(4.0, 4.0)
    dc.DrawRectangle(0.0, 0.0, float64(ctx.Maxx), float64(ctx.Maxy))
    dc.Stroke()
    dcimg := dc.Image()
    bounds := dcimg.Bounds()
    img := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
    draw.Draw(img, img.Bounds(), dcimg, bounds.Min, draw.Src)
    buildTexture(img)
}

func buildTexture(img *image.NRGBA) {
    gl.TexImage2D(
        gl.TEXTURE_2D,
        0,
        gl.RGBA,
        int32(img.Rect.Size().X),
        int32(img.Rect.Size().Y),
        0,
        gl.RGBA,
        gl.UNSIGNED_BYTE,
        gl.Ptr(img.Pix))
}

func generateBuffers() Mesh {
    var mesh Mesh
    mesh.nr_vertices = 6
    var vertexes = []float32{
        -1.0, -1.0, 0.0,
        1.0, -1.0, 0.0,
        -1.0, 1.0, 0.0,
        1.0, -1.0, 0.0,
        -1.0, 1.0, 0.0,
        1.0, 1.0, 0.0,
    }
    var uvs = []float32{
        0.0, 0.0,
        1.0, 0.0,
        0.0, 1.0,
        1.0, 0.0,
        0.0, 1.0,
        1.0, 1.0,
    }
    mesh.program_id = basicShader()

    gl.UseProgram(mesh.program_id)
    gl.GenBuffers(1, &mesh.vertex_id)
    gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vertex_id)
    gl.BufferData(gl.ARRAY_BUFFER, int(4 * (mesh.nr_vertices * 3)), gl.Ptr(vertexes), gl.STATIC_DRAW)
    b_vertex := uint32(gl.GetAttribLocation(mesh.program_id, gl.Str("b_vertex\x00")))
    gl.EnableVertexAttribArray(b_vertex)
    gl.VertexAttribPointer(b_vertex, 3, gl.FLOAT, false, 0, nil)

    gl.GenBuffers(1, &mesh.uv_id)
    gl.BindBuffer(gl.ARRAY_BUFFER, mesh.uv_id)
    gl.BufferData(gl.ARRAY_BUFFER, int(4 * (mesh.nr_vertices * 2)), gl.Ptr(uvs), gl.STATIC_DRAW)
    b_uv := uint32(gl.GetAttribLocation(mesh.program_id, gl.Str("b_uv\x00")))
    gl.EnableVertexAttribArray(b_uv)
    gl.VertexAttribPointer(b_uv, 2, gl.FLOAT, false, 0, nil)

    gl.GenTextures(1, &mesh.tex_id)
    gl.ActiveTexture(gl.TEXTURE0)
    gl.BindTexture(gl.TEXTURE_2D, mesh.tex_id)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
    gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
    u_tex := int32(gl.GetUniformLocation(mesh.program_id, gl.Str("u_tex\x00")))
    gl.Uniform1i(u_tex, 0)

    return mesh
}

func initGL() {
    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)
    gl.ClearColor(1.0, 1.0, 1.0, 1.0)
}

func basicShader() uint32 {
    prog := gl.CreateProgram()

    vertexShaderSource := `
    #version 120

    attribute vec3 b_vertex;
    attribute vec2 b_uv;
    varying vec2 v_uv;

    void main() {
        v_uv = b_uv;
        gl_Position = vec4(b_vertex, 1.0);
    }
` + "\x00"

    fragmentShaderSource := `
    #version 120

    uniform sampler2D u_tex;
    varying vec2 v_uv;

    void main() {
        gl_FragColor = texture2D(u_tex, v_uv);
    }
` + "\x00"

    vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
    if err != nil {
        panic(err)
    }
    fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
    if err != nil {
        panic(err)
    }
    gl.AttachShader(prog, vertexShader)
    gl.AttachShader(prog, fragmentShader)
    gl.LinkProgram(prog)

    return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
    shader := gl.CreateShader(shaderType)

    csources, free := gl.Strs(source)
    gl.ShaderSource(shader, 1, csources, nil)
    free()
    gl.CompileShader(shader)

    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

        return 0, fmt.Errorf("failed to compile %v: %v", source, log)
    }

    return shader, nil
}
