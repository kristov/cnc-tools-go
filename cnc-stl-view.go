package main

import "os"
import "flag"
import "fmt"
import "strings"
import "runtime"
import "github.com/hschendel/stl"
import "github.com/go-gl/glfw/v3.2/glfw"
import "github.com/go-gl/gl/v2.1/gl"
import "github.com/go-gl/mathgl/mgl32"

type Mesh struct {
    vao uint32
    nr_vertices uint32
}

func main() {
    runtime.LockOSThread()

    var stlFile string
    var winWidth int
    var winHeight int

    flag.StringVar(&stlFile, "stl", "", "STL file to view")
    flag.IntVar(&winWidth, "width", 640, "Window width")
    flag.IntVar(&winHeight, "height", 480, "Window height")
    flag.Parse()

    fmt.Fprintln(os.Stderr, "opening:", stlFile)

    solid, err := stl.ReadFile(stlFile)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    err = glfw.Init()
    if err != nil {
        panic(err)
    }
    defer glfw.Terminate()

    window, err := glfw.CreateWindow(winWidth, winHeight, stlFile, nil, nil)
    if err != nil {
        panic(err)
    }

    window.MakeContextCurrent()

    err = gl.Init();
    if err != nil {
        panic(err)
    }

    version := gl.GoStr(gl.GetString(gl.VERSION))
    fmt.Fprintln(os.Stderr, "OpenGL version", version)

    prog := basicShader()
    mesh := makeMesh(solid, prog)

    projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(winWidth)/float32(winHeight), 0.1, 100.0)
    view := mgl32.Ident4()
    model := mgl32.Ident4()

    var mv mgl32.Mat4
    var mvp mgl32.Mat4

    var Z float32 = -30
    model = mgl32.Translate3D(-5, -5, -5)
    view = mgl32.Translate3D(0, 0, Z)
    mv = view.Mul4(model)
    mvp = projection.Mul4(mv)

    window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
        xangle := (float32(xpos) / float32(winWidth)) * 6.2831
        yangle := (float32(ypos) / float32(winHeight)) * 6.2831
        model = mgl32.Translate3D(-5, -5, -5)
        model = mgl32.HomogRotate3DY(-xangle).Mul4(model)
        model = mgl32.HomogRotate3DX(yangle).Mul4(model)
        mv = view.Mul4(model)
        mvp = projection.Mul4(mv)
    })
    window.SetScrollCallback(func(w *glfw.Window, xoff float64, yoff float64) {
        Z = Z + (0.5 * float32(yoff))
        view = mgl32.Translate3D(0, 0, Z)
        mv = view.Mul4(model)
        mvp = projection.Mul4(mv)
    })

    m_mvp_id := gl.GetUniformLocation(prog, gl.Str("m_mvp\x00"))
    m_mv_id := gl.GetUniformLocation(prog, gl.Str("m_mv\x00"))

    for (!window.ShouldClose()) {
        gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
        gl.UseProgram(prog)

        gl.UniformMatrix4fv(m_mvp_id, 1, false, &mvp[0])
        gl.UniformMatrix4fv(m_mv_id, 1, false, &mv[0])

        gl.BindVertexArray(mesh.vao)
        gl.DrawArrays(gl.TRIANGLES, 0, int32(mesh.nr_vertices))

        glfw.PollEvents()
        window.SwapBuffers()
    }
}

func makeMesh(solid *stl.Solid, prog uint32) Mesh {
    var mesh Mesh

    var nr_vertices uint32
    var normals []float32
    var vertices []float32

    triangles := solid.Triangles

    nr_vertices = 0
    for _, triangle := range triangles {
        for _, vertex := range triangle.Vertices {
            normals = append(normals, triangle.Normal[0])
            normals = append(normals, triangle.Normal[1])
            normals = append(normals, triangle.Normal[2])
            vertices = append(vertices, vertex[0])
            vertices = append(vertices, vertex[1])
            vertices = append(vertices, vertex[2])
            nr_vertices++
        }
    }

    var vao uint32
    gl.GenVertexArrays(1, &vao)
    gl.BindVertexArray(vao)

    var vbo uint32
    gl.GenBuffers(1, &vbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
    gl.BufferData(gl.ARRAY_BUFFER, int(4 * (nr_vertices * 3)), gl.Ptr(vertices), gl.STATIC_DRAW)

    b_vertex := uint32(gl.GetAttribLocation(prog, gl.Str("b_vertex\x00")))
    gl.EnableVertexAttribArray(b_vertex)
    gl.VertexAttribPointer(b_vertex, 3, gl.FLOAT, false, 0, nil)

    var nbo uint32
    gl.GenBuffers(1, &nbo)
    gl.BindBuffer(gl.ARRAY_BUFFER, nbo)
    gl.BufferData(gl.ARRAY_BUFFER, int(4 * (nr_vertices * 3)), gl.Ptr(normals), gl.STATIC_DRAW)

    b_normal := uint32(gl.GetAttribLocation(prog, gl.Str("b_normal\x00")))
    gl.EnableVertexAttribArray(b_normal)
    gl.VertexAttribPointer(b_normal, 3, gl.FLOAT, false, 0, nil)

    mesh.vao = vao
    mesh.nr_vertices = nr_vertices

    gl.Enable(gl.DEPTH_TEST)
    gl.DepthFunc(gl.LESS)
    gl.ClearColor(1.0, 1.0, 1.0, 1.0)

    return mesh
}

func basicShader() uint32 {
    prog := gl.CreateProgram()

    vertexShaderSource := `
    #version 120

    attribute vec3 b_vertex;
    attribute vec3 b_normal;

    uniform mat4 m_mvp;

    varying vec3 v_vertex;
    varying vec3 v_normal;

    void main() {
        v_normal = b_normal;
        v_vertex = b_vertex;
        gl_Position = m_mvp * vec4(b_vertex, 1.0);
    }
` + "\x00"

    fragmentShaderSource := `
    #version 120

    uniform mat4 m_mv;

    varying vec3 v_normal;
    varying vec3 v_vertex;

    void main() {
        vec3 normal_ms = normalize(vec3(m_mv * vec4(v_normal, 0.0)));
        vec3 light_ms = vec3(0.0, 0.0, 0.0);
        vec3 vert_ms = vec3(m_mv * vec4(v_vertex, 1.0));
        vec3 stl = light_ms - vert_ms;

        float brightness = dot(normal_ms, stl) / (length(stl) * length(normal_ms));
        brightness = clamp(brightness, 0.0, 1.0);

        vec3 d_color = vec3(0, 1, 0) * brightness;
        gl_FragColor = vec4(d_color, 1.0);
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
