<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    {{- if exists "title" }}
        <title>{{- include "title" }}</title>
    {{- else }}
        <title>Welcome</title>
    {{- end }}
    <style>
      html,body{ min-width: 100%; min-height: 100%; margin: 0; padding: 0; }
      body{ display: block;font-family: 'Courier New', Courier,
      monospace; }
    </style>
    {{- include "styles" }}
    {{- include "scripts" }}
</head>
<body>
    {{- include "app-header" . }}
    {{- view }}
    {{- include "@partials/footer" . }}
</body>
</html>