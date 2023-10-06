# Red Hat Summit Connect 2023 - Automation track

Welcome to the repository for the Automation Track of Red Hat Summit Connect 2023!

In this repository you will find the instructions and configuration to replicate the use cases in your environment.

- [Use cases](#use-cases)
  - [Creation and provisioning of a Virtual Machine in Openshift Virtualization](#creation-and-provisioning-of-a-virtual-machine-in-openshift-virtualization)
    - [VM Creation](#vm-creation)
    - [Day 2 Operations](#day-2-operations)
    - [Automatic Red Hat Insights advisory remediation](#automatic-red-hat-insights-advisory-remediation)
  - [Remediating an alert from OCP Virtualization VM with Service Now and EDA](#remediating-an-alert-from-ocp-virtualization-vm-with-service-now-and-eda)
- [Configuration](#configuration)
- [Requirements](#requirements)
  - [Red Hat Openshift Container Platform](#red-hat-openshift-container-platform)
    - [Required operators](#required-operators)
    - [Network Configuration](#network-configuration)
    - [Storage Configuration](#storage-configuration)
    - [Monitoring Configuration](#monitoring-configuration)
  - [Red Hat Ansible Automation Platform](#red-hat-ansible-automation-platform)

## Use cases

### Creation and provisioning of a Virtual Machine in Openshift Virtualization

The configuration will create a template, **[EDA][OCP] Create VM and configure monitoring** that will provide a cloud-like experience to create a virtual machine on Openshift Virtualization.

You need to fill _Provisioning Webhook_ and _Source URL_ to match your environment configuration.

The use case articulates in three phases:

#### VM Creation

The VM is created instanciating a Virtual Machine template in OCP, the first configuration steps are taken using cloud-init and at the end of the configuration the _Provisioning Webhook_ is called to trigger EDA and start the day 2 operations

#### Day 2 Operations

Once the VM provisioning webhook has been called, the VM is:

- registered on the Red Hat Network to attach a subscriptio
- configured to export metrics using node-exporter for Prometheus
- attached to Red Hat Insights for CVE detection and remediation

#### Automatic Red Hat Insights advisory remediation

Once registered to Red Hat Insights, the platform will send an event that is handled by the **eda-insights** rulebook activation in EDA.

The automation that is triggered will take care of:

- Check if the advisories contain attached security fixes (CVEs)
- Check which of the advisories that contain CVEs need a reboot
- Generate remediation playbooks on Red Hat Insights
- Generate a workflow containing the remediation playbooks
- Attach an approval step to the workflow if the operation requires a reboot

### Remediating an alert from OCP Virtualization VM with Service Now and EDA

## Configuration

## Requirements

In this section you will find the requirements to successfully run all the use cases.

### Red Hat Openshift Container Platform

A working OCP cluster is required to run the demo, as the use case will leverage Red Hat Openshift Virtualization to create a Virtual Machine.

#### Required operators

You will need:

- Openshift Virtualization
- Kubernetes NMState Operator - needed to use bridged network to reach the VMs

#### Network Configuration

The Virtual Machine expects a bridged network using an additional NIC on the Openshift nodes for DHCP and IP reachability from Ansible Controller.

Example files are provided in the [ocp-utils/virt-bridged-network folder](./ocp-utils/virt-bridged-network/)

Ensure to adjust the name of the interface!

#### Storage Configuration

The VM will be configured using two disks, that will consume storage from a default storage class, ensure to have one before proceeding.

#### Monitoring Configuration

One use case relies on OCP Monitoring.
VM-related resources are created by the provisioning playbooks, but some additional steps should be taken:

- Configure monitoring for User Defined workloads
- Create an AlertManager receiver

Configuration snippets can be found in the [ocp-utils/virt-monitoring folder](./ocp-utils/virt-monitoring/)

### Red Hat Ansible Automation Platform
