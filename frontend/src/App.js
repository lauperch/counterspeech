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
          <h1>#noHate</h1>
          <Display/>
        </div>
      </div>
    )
  }
}

export default App