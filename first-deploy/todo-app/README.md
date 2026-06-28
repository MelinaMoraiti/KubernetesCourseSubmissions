**Build the image:** docker build -t todo-app .

**Import the todo-app image with this command:** k3d image import todo-app

**Deploy to Kubernetes Cluster:** kubectl create deployment todoapp-dep --image=todoapp
The deployment's imagePullPolicy was edited after creation with this command kubectl edit deployment todoapp-dep, from always to IfNotPresent.
