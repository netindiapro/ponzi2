#version 450 core

layout(location = 9) uniform sampler2D texture;
layout(location = 10) uniform vec3 mixColor;
layout(location = 11) uniform float mixAmount;

in vec2 texCoord;
in vec3 lighting;

out vec4 fragColor;

void main(void) {
	vec4 color = texture2D(texture, texCoord);
	color = vec4(mix(color.rgb, mixColor, mixAmount), color.a);
	fragColor = vec4(color.rgb * lighting, color.a);
}