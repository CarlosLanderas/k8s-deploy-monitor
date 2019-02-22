import React, { useState, useEffect, useRef } from 'react';
import  deploymentsClient from "./deploymentsClient";
import k8sImage from "./k8s.png";

const sortDeployments = (a, b) =>
  (a.metadata.name > b.metadata.name) ? 1 : ((b.metadata.name > a.metadata.name) ? -1 : 0)

const Deployments = () => {

  const [deploymentsState, setDeployments] = useState([]);
  const refState = useRef();

  const onDeploymentChanged = deployment => {
    
    let deploymentsState = refState.current;
    let prevDeployment = deploymentsState.find(d => d.metadata.name === deployment.metadata.name)
    let deploys;

    if (prevDeployment) {
      prevDeployment.spec.replicas = deployment.spec.replicas;
      deploys = [...deploymentsState.filter(d => d.metadata.name != deployment.metadata.name), prevDeployment]

    } else {
      deploys = [...deploymentsState, deployment]
    }

    deploys.sort(sortDeployments);
    setDeployments(deploys);
  };

  useEffect(async() => {     

      var response = await fetch("/deployments");
      var data = await response.json();
      deploymentsClient(onDeploymentChanged);
      setDeployments(data.items);

  },[]);
  
  useEffect(()=>{
    refState.current = deploymentsState;
  });

  return (
    <div>
      {deploymentsState.map( (deployment,i) =>
      <div key={i}>
        <p>{deployment.metadata.name} has {deployment.spec.replicas} replicas</p>        
        {Array.from(Array(deployment.spec.replicas)).map(i =>  <img height="40" width="40" src={k8sImage}/>)}               
      </div>
      )}      
    </div>
  )
}

export default Deployments;