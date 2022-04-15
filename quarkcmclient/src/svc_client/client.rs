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

// use svc_client::quark_cm_service_client::QuarkCmServiceClient;

use svc_client::quark_cm_service_client::QuarkCmServiceClient;
use svc_client::MaxResourceVersionMessage;
use svc_client::TestRequestMessage;
use tonic::Request;

pub mod svc_client {
    tonic::include_proto!("quarkcmsvc");
}

pub struct Client {
    //client: svc_client::quark_cm_service_client::QuarkCmServiceClient<tonic::transport::Channel>,
}

impl Client {
    pub async fn init() -> Result<(), Box<dyn std::error::Error>> {
        let mut client = QuarkCmServiceClient::connect("http://[::1]:51051").await?;

        let mut name = String::new();
        match hostname::get()?.into_string() {
            Ok(n) => name = n,
            _ => {}
        };
        let request = tonic::Request::new(TestRequestMessage { client_name: name });
        let response = client.test_ping(request).await?;
        println!("TestPing: {:?}", response);

        let all_pods = client.list_pod(()).await?;
        println!("All pods: {:?}", all_pods);

        let all_nodes = client.list_node(()).await?;
        println!("All nodes: {:?}", all_nodes);

        let mut pod_stream = client
            .watch_pod(Request::new(MaxResourceVersionMessage {
                max_resource_version: 0,
            }))
            .await?
            .into_inner();

        while let Some(pod_message) = pod_stream.message().await? {
            println!("Received Pod {:?}", pod_message);
        }

        let mut node_stream = client
            .watch_node(Request::new(MaxResourceVersionMessage {
                max_resource_version: 0,
            }))
            .await?
            .into_inner();

        while let Some(node_message) = node_stream.message().await? {
            println!("Received Node {:?}", node_message);
        }
        Ok(())
    }
}
