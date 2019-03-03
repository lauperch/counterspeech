import React from 'react'
import { Switch, Route } from 'react-router-dom'
import Report from './Report'
import Classifier from './Classifier'

const Display = () => (
  <main>
    <Switch>
      <Route exact path='/' component={Report}/>
      <Route path='/classify' component={Classifier}/>
    </Switch>
  </main>
)

export default Display
