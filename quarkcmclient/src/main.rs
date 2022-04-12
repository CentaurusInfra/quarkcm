use quarkcmsvc::quark_cm_service_client::QuarkCmServiceClient;
use quarkcmsvc::TestRequestMessage;

pub mod quarkcmsvc {
    tonic::include_proto!("quarkcmsvc");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let mut client = QuarkCmServiceClient::connect("http://[::1]:51051").await?;

    let mut name=String::new();
    match hostname::get()?.into_string(){
        Ok(n)=>{name=n},
        _=>{}
    };
    let request=tonic::Request::new(TestRequestMessage{
        client_name: name,
    });
    let response = client.test_ping(request).await?;

    println!("{:?}", response);
    Ok(())
}