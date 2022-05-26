package main

import (
    "os"
    "flag"
    "fmt"
    "strings"
    "runtime"
    "github.com/go-gl/glfw/v3.2/glfw"
    "github.com/go-gl/gl/v2.1/gl"
    "github.com/go-gl/mathgl/mgl32"
    "github.com/fogleman/gg"
    "image"
    "image/draw"
    "io/ioutil"
    "github.com/Succo/wkttoorb"
    "github.com/paulmach/orb"

//    "reflect"
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

    var wktFile string
    var width int
    var height int

    flag.StringVar(&wktFile, "wkt", "", "WKT file to view")
    flag.IntVar(&width, "width", 640, "Window width")
    flag.IntVar(&height, "height", 480, "Window height")
    flag.Parse()

    wktBytes, ferr := ioutil.ReadFile(wktFile)
    if ferr != nil {
        panic(ferr)
    }

    poly, werr := wkttoorb.Scan(string(wktBytes))
    if werr != nil {
        panic(werr)
    }
    //fmt.Println(reflect.TypeOf(geo))
    //fmt.Println(geo)

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
    mesh := generateBuffers(width, height)

    projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(width)/float32(height), 0.1, 100.0)
    view := mgl32.Ident4()
    model := mgl32.Ident4()

    var mv mgl32.Mat4
    var mvp mgl32.Mat4

    var Z float32 = -2
    model = mgl32.Translate3D(0.0, 0.0, 0.0)
    view = mgl32.Translate3D(0.0, 0.0, Z)
    mv = view.Mul4(model)
    mvp = projection.Mul4(mv)

//    window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
//        xangle := (float32(xpos) / float32(width)) * 6.2831
//        yangle := (float32(ypos) / float32(height)) * 6.2831
//        model = mgl32.Translate3D(-5, -5, -5)
//        model = mgl32.HomogRotate3DY(-xangle).Mul4(model)
//        model = mgl32.HomogRotate3DX(yangle).Mul4(model)
//        mv = view.Mul4(model)
//        mvp = projection.Mul4(mv)
//    })
    var scale float64 = 1.0
    window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
        if yoff > 0 {
            scale += 0.1
            draw2D(poly.(orb.Polygon), width, height, scale)
        }
        if yoff < 0 {
            scale -= 0.1
            draw2D(poly.(orb.Polygon), width, height, scale)
        }
    })

    m_mvp_id := gl.GetUniformLocation(mesh.program_id, gl.Str("m_mvp\x00"))
    m_mv_id := gl.GetUniformLocation(mesh.program_id, gl.Str("m_mv\x00"))

    draw2D(poly.(orb.Polygon), width, height, scale)
    for (!window.ShouldClose()) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
        gl.UniformMatrix4fv(m_mvp_id, 1, false, &mvp[0])
        gl.UniformMatrix4fv(m_mv_id, 1, false, &mv[0])
        gl.DrawArrays(gl.TRIANGLES, 0, int32(mesh.nr_vertices))
        glfw.PollEvents()
        window.SwapBuffers()
    }
}

func draw2D(poly orb.Polygon, width int, height int, scale float64) {
    dc := gg.NewContext(width, height)
    dc.ScaleAbout(scale, scale, 0, 0)
    dc.SetRGB(1.0, 0, 0)
    for i := 0; i < len(poly); i++ {
        dc.MoveTo(poly[i][0][0], poly[i][0][1])
        for j := 1; j < len(poly[i]); j++ {
            dc.LineTo(poly[i][j][0], poly[i][j][1])
        }
    }
    dc.Stroke()
    dcimg := dc.Image()
    bounds := dcimg.Bounds()
    img := image.NewNRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
    draw.Draw(img, img.Bounds(), dcimg, bounds.Min, draw.Src)
    rebuildTexture(img)
}

func rebuildTexture(img *image.NRGBA) {
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
    aspect := float32(height) / float32(width)
    var vertexes = []float32{
        -1.0, -aspect, 0.0,
        1.0, -aspect, 0.0,
        -1.0, aspect, 0.0,
        1.0, -aspect, 0.0,
        -1.0, aspect, 0.0,
        1.0, aspect, 0.0,
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

    uniform mat4 m_mvp;
    attribute vec3 b_vertex;
    attribute vec2 b_uv;
    varying vec2 v_uv;

    void main() {
        v_uv = b_uv;
        gl_Position = m_mvp * vec4(b_vertex, 1.0);
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
