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
    "github.com/Succo/wkttoorb"
    "github.com/paulmach/orb"
)

type Mesh struct {
    program_id uint32
    vertex_id uint32
    uv_id uint32
    tex_id uint32
    nr_vertices uint32
    scale float64
}

func main() {
    runtime.LockOSThread()

    var width int
    var height int
    var scale float64 = 1.0

    var things []interface{}
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

    flag.IntVar(&width, "width", 640, "Window width")
    flag.IntVar(&height, "height", 480, "Window height")
    flag.Parse()

    if err := glfw.Init(); err != nil {
        panic(err)
    }
    defer glfw.Terminate()

    window, err := glfw.CreateWindow(width, height, "test", nil, nil)
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
    gl.Viewport(0, 0, int32(width), int32(height))
    mesh := generateBuffers(width, height)
    drawThings(things, width, height, scale)

    window.SetSizeCallback(func(w *glfw.Window, nw int, nh int) {
        width = nw
        height = nh
        gl.Viewport(0, 0, int32(width), int32(height))
        drawThings(things, width, height, scale)
    })

    window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
        if yoff > 0 {
            scale += 0.1
            drawThings(things, width, height, scale)
        }
        if yoff < 0 {
            scale -= 0.1
            drawThings(things, width, height, scale)
        }
    })

    for (!window.ShouldClose()) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
        gl.DrawArrays(gl.TRIANGLES, 0, int32(mesh.nr_vertices))
        glfw.PollEvents()
        window.SwapBuffers()
    }
}

func drawThings(things []interface{}, width int, height int, scale float64) {
    dc := gg.NewContext(width, height)
    dc.ScaleAbout(scale, scale, 0, 0)
    var colors = [][]float64{
        {1.0, 0.0, 0.0},
        {0.0, 1.0, 0.0},
        {0.0, 0.0, 1.0},
    }
    var c uint32 = 0
    dc.SetRGB(colors[c][0], colors[c][1], colors[c][2])
    for i := 0; i < len(things); i++ {
        switch t := things[i].(type) {
            case orb.Polygon:
                drawPolygon(dc, t)
            case orb.LineString:
                drawLineString(dc, t)
            default:
                fmt.Printf("skipping object of unknown type %T\n", t)
        }
        c++
        if c > 2 {
            c = 0
        }
        dc.SetRGB(colors[c][0], colors[c][1], colors[c][2])
    }
    dcimg := dc.Image()
    bounds := dcimg.Bounds()
    img := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
    draw.Draw(img, img.Bounds(), dcimg, bounds.Min, draw.Src)
    buildTexture(img)
}

func drawLineString(dc *gg.Context, ls orb.LineString) {
    dc.MoveTo(ls[0][0], ls[0][1])
    for i := 1; i < len(ls); i++ {
        dc.LineTo(ls[i][0], ls[i][1])
    }
    dc.Stroke()
}

func drawPolygon(dc *gg.Context, poly orb.Polygon) {
    for i := 0; i < len(poly); i++ {
        dc.MoveTo(poly[i][0][0], poly[i][0][1])
        for j := 1; j < len(poly[i]); j++ {
            dc.LineTo(poly[i][j][0], poly[i][j][1])
        }
    }
    dc.Stroke()
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

func generateBuffers(width int, height int) Mesh {
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
