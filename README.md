# Exercises

## Chapter 2: Kubernetes Basics 

- **Exercise 1.1:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.1)
  
- **Exercise 1.2:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.2)
  
- **Exercise 1.3:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.3)

- **Exercise 1.4:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.4)

- **Exercise 1.5:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.5)
  
- **Exercise 1.6:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.6)
  
- **Exercise 1.7:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.7)
  
- **Exercise 1.8:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.8)

- **Exercise 1.9:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.9)
  
- **Exercise 1.10:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.10)

- **Exercise 1.11:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.11)

- **Exercise 1.12:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.12)
  
- **Exercise 1.13:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/1.13)

## Chapter 3: More building Blocks

- **Exercise 2.1:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/2.1)
  
- **Exercise 2.2:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/2.2)

- **Exercise 2.3:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/2.3)

- **Exercise 2.4:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/2.4)
  
- **Exercise 2.5:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/2.5)

- **Exercise 2.6:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/2.6)

- **Exercise 2.7:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/2.7)

- **Exercise 2.8:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/2.8)

- **Exercise 2.9:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/2.9)

- **Exercise 2.10:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/2.10)

## Chapter 4: To the cloud

- **Exercise 3.1:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/3.1)

- **Exercise 3.2:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/3.2)

- **Exercise 3.3:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/3.3)

- **Exercise 3.4:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/3.4)

- **Exercise 3.5:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/3.5)

- **Exercise 3.6:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/3.6)

- **Exercise 3.7:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/3.7)

- **Exercise 3.8:** [Release Link](https://github.com/MelinaMoraiti/KubernetesCourseSubmissions/tree/3.8)

# DBaaS vs Self-Managed Database on GKE (PVC)

Currently runs a self managed Redis. We compare this solution with a self-managed service such as Google Memorystore. 

## Pros/Cons Comparison

**1. Initial Setup**

Self-managed DB: Needs more work, in general. In our case Redis is already set-up so there is nothing to do.

DBaaS: Needs less work in general. With Google Memorystore you can deploy instances in a few clicks from the console, the CLI or client libraries. You then have to repoint your connection strings.

**2. Costs**

Self-managed DB: The infrastructure is cheaper.

DBaaS: Cost is generally higher depends on the configuration

**3. Ongoing Maintenance**

Self-managed DB: Maintenance is more difficult, because is managed by you.

DBaaS: Maintenance is easier, because is managed by the provider.

**4. Backups and Restore**

Self-managed DB: Backups are capable but require manual setup and tests.

DBaaS: For most database solutions backups are built-in and offer simple recovery mechanisms. Cloud SQL backups are automated by default and restores are a command like gcloud sql backups restore, and Memorystore can automate backups
