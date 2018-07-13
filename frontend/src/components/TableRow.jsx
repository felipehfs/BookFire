import React from 'react'

export default (props) => 
    <tr>
        <td>{props.name}</td>
        <td>{props.email}</td>
        <td>{props.description}</td>
        <td>
            <p><a href="#" data-id={props.id} onClick={props.onDelete}>Remover</a> <a href="#" data-id={props.id} onClick={props.onUpdate} index={props.index}>Update</a></p>
        </td>
    </tr>