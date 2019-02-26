import React, { Component } from 'react'

class TodoElements extends Component {
  createTasks = element => {
    return (
      <li key={element.key} onClick={() => this.props.deleteElement(element.key)}>
        {element.text}
      </li>
    )
  }
  render() {
    const todoEntries = this.props.entries
    const listElements = todoEntries.map(this.createTasks)

    return <ul className="theList">{listElements}</ul>
  }
}

export default TodoElements