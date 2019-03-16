import { Component } from 'react';
import React from 'react';


class Classifier extends Component {
  
  constructor() {
    super()
    this.state = {
      currentElement: {
        content: '',
        key: '',
        source: '',
      },
      currentRadio: {
        value: '',
        isHs: false,
        isNotHs: false,
        idk: false
      }
    }
  }

  handleInput = e => {
    const elementText = e.target.value
    const currentElement = { text: elementText, key: Date.now() }
    this.setState({
      currentElement,
    })
  }

  handleSubmit = e => {
    e.preventDefault();
    // sendtoapi

    let url = ""
    if (process.env.NODE_ENV === 'prod') {
      url = 'http://35.198.123.101:5000/submit';
    } else {
      url = 'http://localhost:5000/submit';
    }

    console.log("currentradio",this.state.currentRadio)
    
    fetch(url, {
      method: 'POST',
      mode:'no-cors',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        Content: this.state.currentElement.content,
        IsHS: this.state.currentRadio.isHs === true,
        IsNotHS: this.state.currentRadio.isNotHs === true,
        Idk: this.state.currentRadio.idk === true
      })
    })


    console.log(this.state.currentElement)
    console.log(this.state.currentRadio)
    this.loadNew()
  }

  handleRadioChange = e => {
    switch(e.target.value) {
      case "hsRadio":
        this.setState({currentRadio: {value: e.target.value, isHs: true}})
        break
      case "notHsRadio":
        this.setState({currentRadio: {value: e.target.value, isNotHs: true}})
        break
      default:
        this.setState({currentRadio: {value: e.target.value, idk: true}})
    }
  }

  loadNew = () => {
    let url = ""
    if (process.env.NODE_ENV === 'prod') {
      url = 'http://35.198.123.101:5000/random';
    } else {
      url = 'http://localhost:5000/random';
    }
    fetch(url)
      .then(response => response.json())
      .then(data => {
        const first = data[0]
        //console.log(first)
        this.setState({ currentElement: first })
      });
  }

  componentDidMount() {
    this.loadNew()
  }

  render() {
    return(
      <div className="reportForm">
        <div className="header">
          <form className="submit" onSubmit={this.handleSubmit}>
            <p>{this.state.currentElement.content}</p>
            { this.state.currentElement.url &&(<a href={this.state.currentElement.url} target="_blank">URL</a>)}
            <label><input type="radio" value="hsRadio" onChange={this.handleRadioChange} checked={this.state.currentRadio.value === "hsRadio"}/> i think it's hatespeech</label>
            <label><input type="radio" value="notHsRadio" onChange={this.handleRadioChange} checked={this.state.currentRadio.value === "notHsRadio"} /> i do not think it's hatespeech</label>
            <label><input type="radio" value="idkRadio" onChange={this.handleRadioChange} checked={this.state.currentRadio.value === "idkRadio"}/> i'm not sure</label>
            <button className="classifyButton" type="submit"> Classify </button>
          </form>
        </div>
      </div>
    )
  }
}

export default Classifier;
