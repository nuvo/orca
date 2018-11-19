package utils

import (
	"log"
	"time"
)

// IsEnvValidWithLoopBackOff validates the state of a namespace with back off loop
func IsEnvValidWithLoopBackOff(name, kubeContext string) (bool, error) {
	envValid := false
	maxAttempts := 30
	for i := 1; i <= maxAttempts; i++ {
		envValid, err := IsEnvValid(name, kubeContext)
		if err != nil {
			return envValid, err
		}
		if envValid {
			return envValid, nil
		}
		if i < maxAttempts {
			log.Printf("environment \"%s\" validation failed, will retry in 30 seconds (attempt %d/%d)", name, i, maxAttempts)
			time.Sleep(30 * time.Second)
		}
	}
	return envValid, nil
}

// IsEnvValid validates the state of a namespace
func IsEnvValid(name, kubeContext string) (bool, error) {
	envValid := true
	envValid, err := validatePods(name, kubeContext, envValid)
	if err != nil {
		return envValid, err
	}

	envValid, err = validateEndpoints(name, kubeContext, envValid)
	if err != nil {
		return envValid, err
	}

	return envValid, nil
}

func validatePods(name, kubeContext string, envValid bool) (bool, error) {
	log.Println("validating pods")
	pods, err := getPods(name, kubeContext)
	if err != nil {
		return envValid, err
	}

	log.Println("validating that all pods are in \"Running\" phase")
	for _, pod := range pods.Items {
		phase := pod.Status.Phase
		if phase == "Running" {
			continue
		}
		log.Printf("pod %s is in phase \"%s\"", pod.Name, phase)
		envValid = false
	}

	log.Println("validating that all containers are in \"Ready\" status")
	for _, pod := range pods.Items {
		statuses := pod.Status.ContainerStatuses
		for _, status := range statuses {
			if status.Ready {
				continue
			}
			log.Printf("container %s/%s is not in \"Ready\" status", pod.Name, status.Name)
			envValid = false
		}
	}

	return envValid, nil
}

func validateEndpoints(name, kubeContext string, envValid bool) (bool, error) {
	log.Println("validating endpoints")
	endpoints, err := getEndpoints(name, kubeContext)
	if err != nil {
		return envValid, err
	}

	log.Println("validating that all endpoints have addresses")
	for _, ep := range endpoints.Items {
		subsets := ep.Subsets
		addresses := 0
		for _, subset := range subsets {
			addresses += len(subset.Addresses)
		}
		if addresses != 0 {
			continue
		}
		log.Printf("endpoint %s has no addresses", ep.Name)
		envValid = false
	}

	return envValid, nil
}
