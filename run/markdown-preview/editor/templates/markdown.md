# Playing with Markdown

This UI allows a user to write Markdown text and preview the rendered HTML.

You may be familiar with this composition workflow from sites such as Github
or Wikipedia.

In practice, this web page does the following:

* On click of the *"Preview Rendered Markdown"* button, browser JavaScript
  lifts the markdown text and sends it to the editor UI's public backend.
* The editor backend sends the text on to a private Render service which
  converts it to HTML.
* The HTML is injected into the web page in the right-side **Preview** area.

## Markdown Background

From **[John Gruber on Daring Fireball](https://daringfireball.net/projects/markdown/)**:

> "Markdown is a text-to-HTML conversion tool for web writers. Markdown allows
> you to write using an easy-to-read, easy-to-write plain text format, then
> convert it to structurally valid XHTML (or HTML)."

You can read more about the [syntax on Wikipedia](https://en.wikipedia.org/wiki/Markdown)