docker_build('caffy-beans-example', '.', dockerfile='Dockerfile')
k8s_yaml('deployment.yaml')
k8s_resource('caffy-beans-example', port_forwards=8080)