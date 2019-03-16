import React, { Component } from 'react'
import './App.css'
import Display from './Display'
import NavbarPage from './NavbarPage'

class App extends Component {
  render() { 
    return (
      <div>
        <NavbarPage/>
        <div className="App">
          <Display/>
        </div>
      </div>
    )
  }
}

export default App
