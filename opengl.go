package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (

	// Co-ordinates relative to the center of the window, which is 0,0
	// Can be execute using gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
	// Normal triangle
	/*triangle = []float32{
		0, 0.5, 0, // top
		-0.5, -0.5, 0, // left
		0.5, -0.5, 0, // right
	}*/

	// Right angled triangle
	/*triangle = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}*/

	// Square from two right angled triangles
	square = []float32{
		-0.5, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		-0.5, 0.5, 0,
		0.5, 0.5, 0,
		0.5, -0.5, 0,
	}

	vertexShaderSource = `
#version 410
in vec3 vp;
void main() {
	gl_Position = vec4(vp, 1.0);
}
` + "\x00"

	fragmentShaderSource = `
#version 410
out vec4 frag_colour;
void main() {
	frag_colour = vec4(1, 1, 1, 1);
}
` + "\x00"
)

func initOpenGL() {
	if err := gl.Init(); err != nil {
		log.Fatalf("Error initialising OpenGL: %s\n")
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Printf("OpenGL version: %s\n", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		log.Fatalf("Error compiling vertexShader: %s\n", err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		log.Fatalf("Error compiling fragmentShader: %s\n", err)
	}

	program = gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

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

func initGlfw() {
	if err := glfw.Init(); err != nil {
		log.Fatalf("Error initialising glfw: %s\n", err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	_window, err := glfw.CreateWindow(WIDTH, HEIGHT, "Conway's Game of Life", nil, nil)
	if err != nil {
		log.Fatalf("Error creating window: %s\n", err)
	}

	// TODO(): Figure out the correct way to do this...
	window = _window
	window.MakeContextCurrent()

	return
}
