import React, { useState, useEffect, useRef } from 'react';
import  deploymentsClient from "./deploymentsClient";
import k8sImage from "./k8s.png";
import k8sRedImage from "./k8s-red.png";

const sortDeployments = (a, b) =>
  (a.metadata.name > b.metadata.name) ? 1 : ((b.metadata.name > a.metadata.name) ? -1 : 0);

const Deployments = () => {

  const [deploymentsState, setDeployments] = useState([]);
  const refState = useRef();

  const onDeploymentChanged = deployment => {
    
    let state = refState.current;
    let prevDeployment = state.find(d => d.metadata.name === deployment.metadata.name);
    let deploys;

    if (prevDeployment) {
      prevDeployment.spec.replicas = deployment.spec.replicas;
      prevDeployment.status.availableReplicas = deployment.status.availableReplicas;
      deploys = [...state.filter(d => d.metadata.name != deployment.metadata.name), prevDeployment];

    } else {
      deploys = [...state, deployment]
    }

    deploys.sort(sortDeployments);    
    setDeployments(deploys);
    
  };

  useEffect(async() => {     

      var response = await fetch("/deployments");
      var data = await response.json();
      setDeployments(data.items);
      
      deploymentsClient(onDeploymentChanged);      
  },[]);
  
  useEffect(()=>{
    refState.current = deploymentsState;       
  });

  const renderAvailableReplicas = (deployment) => {
      let available = [];
      for(var i = 0 ; i < (deployment.status.availableReplicas || 0); i++) {
        available.push(<img key={i} className="k8simage" src={k8sImage}/>);
      }
      return available;
  };
  
  const renderUnavailableReplicas = (deployment) => {
    let unavailable = [];
    for(var i = 0; i < deployment.spec.replicas - (deployment.status.availableReplicas || 0); i++) {
      unavailable.push(<img key={i} className="k8simage" src={k8sRedImage}/>);
    }
    return unavailable;
  };

  return (
    <div className="deployments">
      {deploymentsState.map( (deployment,i) =>
      <div key={i}>
        <p>{deployment.metadata.name} has {deployment.spec.replicas} replicas ( {deployment.status.availableReplicas || 0 } available and { (deployment.spec.replicas || 0 ) - (deployment.status.availableReplicas || 0 )} not available) </p>        
        {renderAvailableReplicas(deployment)}
        {renderUnavailableReplicas(deployment)}        
                      
      </div>
      )}      
    </div>
  )
}

export default Deployments;