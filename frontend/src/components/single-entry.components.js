import React from "react";
import 'bootstrap/dist/css/bootstrap.css';

import {Button, Card, Row, Col} from "react-bootstrap";

const Entry = ({entryData, setChangeIngredient, deleteSingleEntry, setChangeEntry}) => {
    return (
        <Card>
            <Row>
                <Col>Dish:{entryData !== undefined && entryData.dish}</Col>
                <Col>Ingredients:{entryData !== undefined && entryData.ingredients}</Col>
                <Col>Calories:{entryData !== undefined && entryData.calories}</Col>
                <Col>Fat:{entryData !== undefined && entryData.fat}</Col>
                <Col><Button onClick={()=> deleteSingleEntry(entryData._id)}>Delete Entry</Button></Col>
                <Col><Button onClick={()=>changeIngredient()}>Change Ingredients</Button></Col>
                <Col><Button onClick={()=>changeEntry()}>Change Entry</Button></Col>
            </Row>
        </Card>
    )
    function changeIngredient(){
        setChangeIngredient(
            {
                "change":true,
                "id": entryData._id
            }
        )
    }
    function changeEntry(){
        setChangeEntry(
            {
                "change":true,
                "id":entryData._id
            }
        )
    }
}

export default Entry