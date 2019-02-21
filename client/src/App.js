import React, { useState, useEffect } from 'react';
import { deploymentsClient } from "./deploymentsClient";
import './App.css';

  
function DeploymentWatcher(){
debugger;
 const [deployments, setDeployment] = useState([]);  
 const onMessage = (deployment) =>  {
   
    setDeployment([...deployments, deployment]);
 }

  useEffect(()=>{
    deploymentsClient(onMessage);
  },[]);

  return(
    <div>
      {deployments.map(deployment => {
        return <span>{deployment.metadata.name} has {deployment.spec.replicas} replicas</span>
      })}
    </div>
  )
}

export default DeploymentWatcher;