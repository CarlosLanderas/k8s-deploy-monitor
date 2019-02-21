import React, { useState, useEffect, useRef } from 'react';
import  deploymentsClient from "./deploymentsClient";


const sortDeployments = (a, b) =>
  (a.metadata.name > b.metadata.name) ? 1 : ((b.metadata.name > a.metadata.name) ? -1 : 0)

const Deployments = () => {
  console.log("Executing init");
  
  const [deploymentsState, setDeployments] = useState([]);
  const refState = useRef();
  
  console.log(deploymentsState);

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
      deploymentsClient(onDeploymentChanged)
      setDeployments(data.items);

  },[]);
  
  useEffect(()=>{
    refState.current = deploymentsState;
  });

  return (
    <div>
      {deploymentsState.map(deployment =>
        <p>{deployment.metadata.name} has {deployment.spec.replicas} replicas</p>
      )}      
    </div>
  )
}

export default Deployments;