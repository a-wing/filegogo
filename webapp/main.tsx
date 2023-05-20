import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './app2'

import 'virtual:windi.css'
//import 'virtual:windi-devtools'
//import 'virtual:windi-base.css'
//import 'virtual:windi-components.css'
//import 'virtual:windi-utilities.css'

import GithubCorner from 'react-github-corner'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App/>
    <GithubCorner href="https://github.com/a-wing/filegogo" bannerColor="#64CEAA" octoColor="#FFF" />
  </React.StrictMode>
)
