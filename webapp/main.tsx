import React from "react"
import ReactDOM from "react-dom/client"
import App from "./app"

import { Provider } from "jotai"

import "virtual:uno.css"
import "@unocss/reset/tailwind.css"

import GithubCorner from "react-github-corner"

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Provider>
      <App/>
    </Provider>
    <GithubCorner href="https://github.com/a-wing/filegogo" bannerColor="#64CEAA" octoColor="#FFF" />
  </React.StrictMode>
)
