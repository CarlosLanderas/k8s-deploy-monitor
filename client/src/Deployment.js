import React, { useState, useEffect, useRef } from 'react';
import  deploymentsClient from "./deploymentsClient";
import k8sImage from "./k8s.png";
import k8sRedImage from "./k8s-red.png";

const sortDeployments = (a, b) =>
  (a.metadata.name > b.metadata.name) ? 1 : ((b.metadata.name > a.metadata.name) ? -1 : 0)

const Deployments = () => {

  const [deploymentsState, setDeployments] = useState([]);
  const [availableReplicas, setAvailableReplicas] = useState([]);
  const [unavailableReplicas, setUnavailableReplicas] = useState([]);
  const refState = useRef();

  const onDeploymentChanged = deployment => {
    
    let deploymentsState = refState.current;
    let prevDeployment = deploymentsState.find(d => d.metadata.name === deployment.metadata.name)
    let deploys;

    if (prevDeployment) {
      prevDeployment.spec.replicas = deployment.spec.replicas;
      deploys = [...deploymentsState.filter(d => d.metadata.name != deployment.metadata.name), Object.assign({},prevDeployment)]

    } else {
      deploys = [...deploymentsState, Object.assign({}, deployment)]
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

  const renderAvailableReplicas = (deployment) => {
      let available = [];
      for(var i = 0 ; i < deployment.status.availableReplicas; i++) {
        available.push(<img key={i} height="40" width="40" src={k8sImage}/>);
      }
      return available;
  };
  
  const renderUnavailableReplicas = (deployment) => {
    let unavailable = [];
    for(var i = 0; i < deployment.spec.replicas - deployment.status.availableReplicas; i++) {
      unavailable.push(<img key={i} height="40" width="40" src={k8sRedImage}/>);
    }
    return unavailable;
  };

  return (
    <div className="deployments">
      {deploymentsState.map( (deployment,i) =>
      <div key={i}>
        <p>{deployment.metadata.name} has {deployment.spec.replicas} replicas ( {deployment.status.availableReplicas} 
         available and { deployment.spec.replicas - deployment.status.availableReplicas} not available) </p>        
        {renderAvailableReplicas(deployment)}
        {renderUnavailableReplicas(deployment)}        
                      
      </div>
      )}      
    </div>
  )
}

export default Deployments;