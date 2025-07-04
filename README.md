# Kuadrant-MCP-Walkthrough
Prerequisites
- A Kubernetes cluster 
- Kuadrant installed, please see the Kuadrant [docs](https://docs.kuadrant.io/latest/install-helm/) for more information.

## Deploy MCP Gateway aka MCP server(s)

1. Create the MCP Server namespace:

    ```sh
    kubectl create ns mcp-server
    ```

1. Deploy the MCP gateway server.

    Note: We have created a image for the MCP Gateway that can be used to deploy the MCP server(s) in your Kubernetes cluster.

    ```sh
    kubectl apply -f mcp-gateway/
    ```
    For more information about the MCP Gateway piece, see the [MCP Gateway repository](https://github.com/david-martin/mcp-gateway-poc) .


1.  Ensure the mcp gateway (Server) is up and running:
    ```sh
    kubectl get pods -n mcp-server
    ```

## Deploying Kuadrant policies 

1. Clone the Repository
    ```sh
    git clone https://github.com/your-org/Kuadrant-MCP-Walkthrough.git
    cd Kuadrant-MCP-Walkthrough
    ```

 1. Create the gateway namespace:

    ```sh
        kubectl create ns mcp-gateway
    ```

1. Create the secret credentials in the same namespace as the Gateway - these will be used to configure DNS:

    ```sh
    kubectl -n mcp-gateway create secret generic aws-credentials \
    --type=kuadrant.io/aws \
    --from-literal=AWS_ACCESS_KEY_ID=$KUADRANT_AWS_ACCESS_KEY_ID \
    --from-literal=AWS_SECRET_ACCESS_KEY=$KUADRANT_AWS_SECRET_ACCESS_KEY
    ```

1. Create the secret credentials in the cert-manager namespace:

    ```sh

    kubectl -n cert-manager create secret generic aws-credentials \
    --type=kuadrant.io/aws \
    --from-literal=AWS_ACCESS_KEY_ID=$KUADRANT_AWS_ACCESS_KEY_ID \
    --from-literal=AWS_SECRET_ACCESS_KEY=$KUADRANT_AWS_SECRET_ACCESS_KEY
    ```   

1. Apply the policies
    ```sh
    kubectl apply -f kuadrant/
    ```
    This will deploy :
    * Gateway
    * AuthN Auth policy
    * AuthZ Auth policy for specific MCP tools
    * RateLimit policy
    * RateLimit policy for specific MCP tool
    * Dns policy
    * HTTPRoute
  

## Setting Up MCP Inspector to test your policies

Note: Adding custom hearers to MCP Inspector is not supported in the current version. Theres is a PR open with this feature implemented, which works perfectly so that can be used instead https://github.com/modelcontextprotocol/inspector/pull/549 


1. Clone the MCP Inspector Repository
    ```bash
    git clone https://github.com/popomore/inspector.git
    cd mcp-inspector
    ```

1. Deploy MCP Inspector
    ```bash
    npm run build
    npm start
    ```

1. Access MCP Inspector
When the MCP Inspector is running, you can access it via the following URL. The auth token will be outputted in the terminal when you run the inspector. Replace `<your-auth-token>` with the actual token provided by the MCP Inspector.
    ```
    http://localhost:6274/?MCP_PROXY_AUTH_TOKEN=<your-auth-token>
    ```
1. From there you can setup the MCP Inspector to use authorisation and test out the tools the MCP server provides.

