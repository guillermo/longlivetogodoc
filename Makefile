

DEPS=128.png 16.png 256.png 32.png 48.png go.dev.css go.dev.js manifest.json

longlivegodoc.zip: ${DEPS}
	zip -r longlivegodoc.zip ${DEPS}

