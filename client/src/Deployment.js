import React, { useState, useEffect } from 'react';
import  deploymentsClient from "./deploymentsClient";


const sortDeployments = (a, b) =>
  (a.metadata.name > b.metadata.name) ? 1 : ((b.metadata.name > a.metadata.name) ? -1 : 0)

let deploymentsSocket = new deploymentsClient();

const Deployments = () => {
  console.log("Executing init");
    
  const [deploymentsState, setDeployments] = useState([]);
  console.log(deploymentsState);

  const onDeploymentChanged = deployment => {
    console.log(deploymentsState);
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

  deploymentsSocket.onMessage(onDeploymentChanged);
 
  useEffect(() => {
      
    (async() => {
      var response = await fetch("/deployments");
      var data = await response.json();
      
      setDeployments(data.items);
    })();

  }, []);

  return (
    <div>
      {deploymentsState.map(deployment =>
        <p>{deployment.metadata.name} has {deployment.spec.replicas} replicas</p>
      )}      
    </div>
  )
}

export default Deployments;