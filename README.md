# Welcome to the gothex template

This template is meant to make scarfolding your App easy. It is a full stack application using Golang's Echo framework, HTMX in the frontend, Templ templating language & tailwind CSS with user session management and browser hot reload.

For this project you will need the following installed on your machine:

✅ Golang => https://go.dev/doc/install<br>
✅ Task => https://taskfile.dev/installation<br>
✅ Gowatch => https://github.com/silenceper/gowatch<br>
✅ Templ => https://templ.guide/quick-start/installation<br>
✅ Tailwind CLI => https://tailwindcss.com/docs/installation<br>

Task is just meant to make execution easy you can use Make if you wish or just run the comand directly from your terminal

Your <strong>.env</strong> file must look something like this:

```php
PORT=3000
SESSION_AGE=3600
SIGNIN_PAGE=/signin
AFTER_SIGNIN_PAGE=/
APP_TITLE=Gothex Template App
GOTHEX_SECRET=6ub7dedx43rg7bm81gp3k8r
GOOGLE_SECRET=
GOOGLE_CALLBACK=http://localhost:3000/api/google/callback
GOOGLE_KEY=
```

> <strong>Important</strong> <p>In this application, instead of using the html/template package (native Golang templates), we use the a-h/templ library. This amazing library implements a templating language (very similar to JSX) that compiles to Go code. Templ will allow us to write code almost identical to Go (with expressions, control flow, if/else, for loops, etc.) and have autocompletion thanks to strong typing. This means that errors appear at compile time and any calls to these templates (which are compiled as Go functions) from the handlers side will always require the correct data, minimizing errors and thus increasing the security and speed of our coding.</p>

On the other hand, we use Golang's Echo web framework, which as stated on its website is a "High performance, extensible, minimalist Go web framework".

The use of </>htmx allows behavior similar to that of a SPA, without page reloads when switching from one route to another or when making requests (via AJAX) to the backend.

The styling of the views is achieved through Tailwind CSS is generated at runtime all configs are in the task file.

Finally, minimal use of Alpine Ja is made to achieve the action of closing the alerts when they are displayed or giving interactivity to the show/hide password button in its corresponding input.

To run this project, just open 3 terminal windows and run the following commands.

```bash
task tailwind
```

```bash
task server
```

```bash
task templ
```
