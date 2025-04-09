<section>
  <h1>Contact Page</h1>
  <p>Contact me at: xxx@yyy.zzz</p>
  {{- include "contact-form" }}
  {{- include "pages/contact/social" }}
</section>
{{- define "scripts"}}
<script>
  alert("Call me!")
</script>
{{- end}}
{{- define "custom-nav"}}Call{{- end}}
{{- define "title"}}Contact Me!{{- end}}