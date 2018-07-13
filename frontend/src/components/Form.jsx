import React from 'react'
import "./Form.css"
export default (props) => 
    <form className="form" onSubmit={props.onSubmit}>
        <label>Name</label><br/>
        <input type='text' name="name" onChange={props.onChange} value={props.userInput.name}/><br/>
        <label>Email</label><br/>
        <input type="email" name="email" onChange={props.onChange} value={props.userInput.email}/><br/>
        <label>Description</label><br/>
        <input type="text" name="description" onChange={props.onChange} value={props.userInput.description}/><br/>
        <button type="submit">Confirmar</button>
    </form>