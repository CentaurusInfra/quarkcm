/*
Copyright 2022 quarkcm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

use super::super::RDMA_CTLINFO;
use crate::rdma_ctrlconn::*;
use svc_client::quark_cm_service_client::QuarkCmServiceClient;
use svc_client::MaxResourceVersionMessage;
use tonic::Request;
use tokio::time::*;

pub mod svc_client {
    tonic::include_proto!("quarkcmsvc");
}

#[derive(Debug)]
pub struct PodInformer {
    pub max_resource_version: i32,
}

impl PodInformer {
    pub async fn new() -> Result<Self, Box<dyn std::error::Error>> {
        let mut informer = Self {
            max_resource_version: 0,
        };
        let mut client = QuarkCmServiceClient::connect("http://[::1]:51051").await?;

        let ref pods_message = client.list_pod(()).await?.into_inner().pods;
        println!("All pods: {:#?}", pods_message);
        if pods_message.len() > 0 {
            let mut pods_map = RDMA_CTLINFO.pods.lock();

            for pod_message in pods_message {
                let pod = Pod {
                    key: pod_message.key.clone(),
                    ip: pod_message.ip.clone(),
                    node_name: pod_message.node_name.clone(),
                    resource_version: pod_message.resource_version,
                };
                pods_map.insert(pod_message.ip.clone(), pod);
                if pod_message.resource_version > informer.max_resource_version {
                    informer.max_resource_version = pod_message.resource_version;
                }
            }
            println!("max_resource_version: {}", informer.max_resource_version);
            println!("RDMA_CTLINFO {:#?}", pods_map);
        }

        Ok(informer)
    }
}

impl PodInformer {
    pub async fn run(&mut self) -> Result<(), Box<dyn std::error::Error>> {
        loop {
            let mut client = QuarkCmServiceClient::connect("http://[::1]:51051").await?;

            let mut pod_stream = client
                .watch_pod(Request::new(MaxResourceVersionMessage {
                    max_resource_version: self.max_resource_version,
                    // max_resource_version: 0,
                }))
                .await?
                .into_inner();
    
            while let Some(pod_message) = pod_stream.message().await? {
                println!("Received Pod {:?}", pod_message);
            }

            if *RDMA_CTLINFO.exiting.lock() {
                break;
            } else {
                println!("Wait 1 second for next iteration of watching pod.");
                sleep(Duration::from_secs(1)).await;
            }
        }
        
        Ok(())
    }
}

impl PodInformer {
    fn handle(&mut self) {
        
    }
}
