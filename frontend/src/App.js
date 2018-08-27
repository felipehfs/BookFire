import React, { Component } from 'react';
import './App.css';
import Header from './components/Header'
import Table from './components/Table'
import Form from './components/Form'
import axios from 'axios'

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzgwMTA0NzUsImlheCI6IjIwMTgtMDgtMjJUMjI6MDc6NTUuMTQzMzc3MzItMDM6MDAiLCJzdWIiOiJhZG1pbiJ9.hBrizG4260n7TI2lPV55rCI9h-xNTnzKB86N4wIKB68"
const header = {headers: {
  "Authorization": `Bearer ${token}`
}}

class App extends Component {
  constructor(props){
    super(props)
    this.state = {
      id: '',
      userInput: {
        name: '',
        email: '',
        description: ''
      },
      artists: []
    }
    this.onClear = this.onClear.bind(this)
  }

  onChange(e){
    const {state} = this 
    let input = e.target.name 
    state.userInput[input] = e.target.value 
    this.setState(state)
  }

  onClear() {
    const {state} = this
    state.id = ''
    Object.keys(state.userInput).forEach(key => state.userInput[key] = "")
    this.setState(state)
  }

  componentDidMount(){
    axios.get(`http://localhost:8081/artists/`, header)
    .then(resp => {
      const data = resp.data === null? []: resp.data 
      this.setState({...this.state, artists: data})
    })
    .catch(e => this.setState({...this.state, artists: []}))
  }

  refresh(){
    axios.get(`http://localhost:8081/artists/`, header)
    .then(resp => {
      const data = resp.data === null? []: resp.data 
      this.setState({...this.state, artists: data})
      this.onClear()
    })
    .catch(e => this.setState({...this.state, artists: []}))
  }

  onSubmit(e) {
    e.preventDefault()
    const {userInput} = this.state 
  
    if (this.state.id === ''){
      axios.post("http://localhost:8081/artists/", userInput, header)
      .then(resp => {
        this.refresh()
      }).catch(e => console.log(e))
    } else {
      const id = this.state.id
      axios.put(`http://localhost:8081/artists/${id}`,userInput, header)
      .then(resp => this.refresh())
    }  
  }
  onDelete(e) {
    e.preventDefault()
    const id = e.target.getAttribute("data-id")
    axios.delete(`http://localhost:8081/artists/${id}`, header)
    .then(resp => this.refresh())
  } 
  onUpdate(e) {
    const index = Number.parseInt(e.target.getAttribute("index"))
    const {state} = this 

    const updated = state.artists[index]
    state.id = updated.id
    state.userInput.name = updated.name 
    state.userInput.email = updated.email 
    state.userInput.description = updated.description

    this.setState(state) 
  }
  render() {
    return (
      <div className="App">
        <Header />
        <Form onChange={this.onChange.bind(this)} onClear={this.onClear} userInput={this.state.userInput}
          onSubmit={this.onSubmit.bind(this)}/>
        <Table lista={this.state.artists} 
          onUpdate={this.onUpdate.bind(this)}
          onDelete={this.onDelete.bind(this)}/>
      </div>
    );
  }
}

export default App;
