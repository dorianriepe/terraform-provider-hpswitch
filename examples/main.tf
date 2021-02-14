terraform {
  required_providers {
    hpswitch = {
      version = "0.1.1"
      source  = "dorianriepe/test/hpswitch"
    }
  }
}

provider "hpswitch" {
    host = "my-switch-url.local"
    username = "user"
    password = "ULTRASECRETP4SSW0RD"
}

resource "hpswitch_vlan" "my_vlan" {
    vlan = "3125"
    description = "My Vlan"
    tagged_ports {
      port = "Ten-GigabitEthernet2/0/7"
    }
    tagged_ports {
      port = "Ten-GigabitEthernet2/0/1"
    }
    tagged_ports {
      port = "Ten-GigabitEthernet2/0/2"
    }
}

output "my_vlan" {
  value = hpswitch_vlan.my_vlan
}