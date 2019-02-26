import React, { Component } from 'react'
import './App.css'
import TodoList from './TodoList'
import TodoElements from './TodoElements'

class App extends Component {
  inputElement = React.createRef()
  
  constructor() {
    super()
    this.state = {
      elements: [],
      currentElement: {
        text: '',
        key: '',
      },
    }
  }

  deleteElement = key => {
    const filteredelements = this.state.elements.filter(element => {
      return element.key !== key
    })
    this.setState({
      elements: filteredelements,
    })
  }

  handleInput = e => {
    const elementText = e.target.value
    const currentElement = { text: elementText, key: Date.now() }
    this.setState({
      currentElement,
    })
  }

  addElement = e => {
    e.preventDefault()
    const newElement = this.state.currentElement
    if (newElement.text !== '') {
      const elements = [...this.state.elements, newElement]
      this.setState({
        elements: elements,
        currentElement: { text: '', key: '' },
      })
    }
  }

  render() {
    return (
      <div className="App">
        <h1>TODO</h1>
        <TodoList
          addElement={this.addElement}
          inputElement={this.inputElement}
          handleInput={this.handleInput}
          currentElement={this.state.currentElement}
        />
        <TodoElements entries={this.state.elements} deleteElement={this.deleteElement} />
      </div>
    )
  }
}

export default App