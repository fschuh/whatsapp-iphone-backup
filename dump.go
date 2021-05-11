package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func (app *App) DumpSessions(css []Session) {
	// index.html
	data := `
<frameset cols="25%,*">
	<frame src="sessions/index.html">
	<frame name="session">
</frameset>	
	`
	filename := path.Join(app.DstDir, "index.html")
	ioutil.WriteFile(filename, []byte(data), 0600)
	// mkdir sessions
	os.Mkdir(path.Join(app.DstDir, "sessions"), 0700)
	// sessions.html
	tpl := `
	<h1>WhatsApp</h1>

	{{range .}}
		<a href="session_{{ .ID}}.html" target="session">{{ .Name}}</a><br>
	{{end}}
	
	`
	t, err := template.New("foo").Parse(tpl)
	check("DumpSessions template parsing", err)

	out, err := os.Create(path.Join(app.DstDir, "sessions", "index.html"))
	check("DumpSessions creating file", err)
	defer out.Close()

	err = t.Execute(out, css)
	check("DumpSessions executing template", err)
}

func (app *App) DumpSession(session Session, messages []Message) {
	tpl := `
<style><!--
body {
	background:rgb(229,221,213);	
}
.chat {
	width:600px;
	margin:auto;
}
.message {
	margin: 0.4em;
	padding: 0.5em;
	min-height: 1.5em;
}
.message img { max-width: 100%; display: block; }
.message video { max-width: 100%; display: block; }
.incoming {
	background:white;
}
.outgoing {
	background:rgb(221,247,200);
	text-align:right;
}
.name {
	display: inline-block;
	color: darkslategray;
	font-size: 1.05em;
}
.date {
    display: inline-block;
    color: darkslategray;
    padding-left: 0.5em;
    font-size: 0.90em;
}
.content {
    display: block;
    padding-top: .5em;
}
--></style>
<h1>WhatsApp</h1>

<div class="chat">
{{range .}}	
	<p class="message {{if .JID}}incoming{{else}}outgoing{{end}}">
	    {{ if and .Name .JID }}
          <span class="name">{{ .Name}}</span>
		{{ end }}
		<span class="date">{{ .Date }}</span>
		<span class="content">
			{{ nl2br .Text }}
			{{ if eq .MediaExt ".jpg" }}
				<img loading="lazy" src="../{{.Media}}">
			{{ else if eq .MediaExt ".png" }}
				<img loading="lazy" src="../{{.Media}}">
			{{ else if eq .MediaExt ".webp" }}
				<img loading="lazy" src="../{{.Media}}">
			{{ else if eq .MediaExt ".mp4" }}
				<video preload="metadata" controls>
					<source src="../{{.Media}}" type="video/mp4">
				</video>
			{{ else if eq .MediaExt ".opus" }}
				<audio preload="none" src="../{{.Media}}" controls />
				Your browser does not support the audio element. Try Firefox.
				</audio>
			{{ else if .MediaExt }}
				<a href="../{{.Media}}" />FILE: {{.Media}}</a>
			{{end}}
		</span>
	</p>
{{end}}
</div><!-- chat -->
	`
	funcs := template.FuncMap{
		"nl2br": func(text string) template.HTML {
			return template.HTML(strings.Replace(template.HTMLEscapeString(text), "\n", "<br>", -1))
		},
	}
	t, err := template.New("foo").Funcs(funcs).Parse(tpl)
	check("DumpSession template parsing", err)

	out, err := os.Create(path.Join(app.DstDir, "sessions", fmt.Sprintf("session_%d.html", session.ID)))
	check("DumpSession creating file", err)
	defer out.Close()

	//optional: only dump last 30 messages per chat to HTML
	//if len(messages) > 30 {
	//	messages = messages[len(messages)-30:]
	//}
	err = t.Execute(out, messages)
	check("DumpSession executing template", err)
}
