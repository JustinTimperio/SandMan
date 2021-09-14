<p align="center">
  <img src="https://upload.wikimedia.org/wikipedia/en/8/8b/Sandman_1974_issue1.jpg">
</p>

# Sandman
SandMan is a set of advanced tools for detecting and evading malware analysis sandboxes. Based on the work of ColdFire and VM-Detection, SandMan combines and expands upon these detection methods into a full toolset for evading sandboxes.

## Why SandMan?

In an effort to provide blue-team members better resources, SandMan is a transparent way to test the most effective detection and evasion techniques being used by Malware today. Additonally, those who build red-team tools will benfit from a fully self-contained evasion and detection module for attack simulations and pen-testing.

## How It Works
SandMan uses a varitey of scoreing methods to return a score which represents the likelyhood that the OS is being run as a VM. 

### Detection 

#### Common Checks
- CPU Check
- Ram Check
- Known Mac Address Check
- Time Compression Check

#### Linux VM Detection
- Check DMI Table for VM Entries
-	See if Kernel Detects a Hypervisor
- Check for Hypervisor Flag or User Mode Linux
- Check the Device Tree for VM artifacts
- Look for VM Tools in Modules

#### Windows VM Detection
- Checks the Registry for Blacklisted Keys and Vendors 
- Checks the Device Tree for VM artifacts

### Evasion

