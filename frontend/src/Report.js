import React, { Component } from 'react'

class Report extends Component {
  inputElement = React.createRef()
  
  constructor() {
    super()
    this.state = {
      currentElement: {
        text: '',
        key: '',
      },
    }
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
    console.log('adding', newElement)

    let url = 'http://35.198.123.101:5000/submit'
    
    fetch(url, {
      method: 'POST',
      mode:'no-cors',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        Content: newElement.text
      })
    })

    if (newElement.text !== '') {
      this.setState({
        currentElement: { text: '', key: '' },
      })
    }
  }

  componentDidUpdate() {
    this.inputElement.current.focus()
  }

  render() {
    return (
      <div className="reportForm">
        <div className="header">
            <form className="submit" onSubmit={this.addElement}>
              <textarea
                placeholder="Insert text"
                ref={this.inputElement}
                value={this.state.currentElement.text}
                onChange={this.handleInput}
              />
              <button type="submit"> Report </button>
            </form>
        </div>
      </div>
    )
  }
}

export default Report;
