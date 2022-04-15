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

pub mod rdma_ctrlconn;

mod svc_client;
//use svc_client::client::Client;
use crate::rdma_ctrlconn::*;
use svc_client::pod_informer::PodInformer;

use lazy_static::lazy_static;

lazy_static! {
    pub static ref RDMA_CTLINFO: CtrlInfo = CtrlInfo::default();
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    //Client::init().await?;
    //PodInformer::run().await?;

    let mut pod_informer = PodInformer::new().await?;
    tokio::join!(
        pod_informer.run(),
        sleep_then_print(1),
        sleep_then_print(2),
    );
    //pod_informer.run().await?;
    println!("ok");
    Ok(())
}

async fn sleep_then_print(timer: i32) {
    println!("Start timer {}.", timer);

    tokio::time::sleep(tokio::time::Duration::from_secs(1)).await;
//                                            ^ execution can be paused here

    println!("Timer {} done.", timer);
}
