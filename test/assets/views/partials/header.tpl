{{define "app-header"}}
<header>
  <nav>
    <a href="/">Home</a>
    <a href="/contact">Contact</a>
    {{- if exists "custom-nav"}}
    <a href="/error">{{include "custom-nav"}}</a>
    {{- end}}
  </nav>
</header>
{{end}}