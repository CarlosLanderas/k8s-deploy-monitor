import React, { useState, useEffect } from 'react';
import { deploymentsClient } from "./deploymentsClient";
import Deployment from "./Deployment";
import './App.css';

export default function App() {
  return (<Deployment/>)
}