import os, sys
from urllib.request import urlretrieve

os.chdir(os.path.dirname(sys.argv[0]))

from diagrams import Cluster, Diagram, Edge
from diagrams.k8s.compute import Pod
from diagrams.onprem.queue import Kafka
from diagrams.onprem.database import Postgresql
from diagrams.generic.compute import Rack  # ou diagrams.generic.blank.Blank

with Diagram("campaign producer api"):
   blueline = Edge(color="blue", style="bold")
   blackline = Edge(color="black", style="bold")
   bidirectional_edge = Edge(color="darkOrange", style="bold", forward=True, reverse=True)

   with Cluster("service"):
      producerAPI = Pod("campaign-producer-api")

   with Cluster("queue"):
      producerkafkaOwner = Kafka("campaign.campaign-owner")
      producerKafkaSlug = Kafka("campaign.campaign-slug")
      producerKafkaRegion = Kafka("campaign.campaign-region")
      producerKafkaMerchant = Kafka("campaign.campaign-merchant")
      producerKafkaCampaign = Kafka("campaign.campaign")
      producerKafkaSpent = Kafka("campaign.campaign-spent")
   
   with Cluster("cache"):
      producerLocalCache = Rack("Local Cache")  # Alteração feita aqui
   
   with Cluster("db"):
      producerDb = Postgresql("campaign-consumer-db")

   # Definindo as conexões
   producerAPI - blueline >> producerkafkaOwner
   producerAPI - blueline >> producerKafkaSlug
   producerAPI - blueline >> producerKafkaRegion
   producerAPI - blueline >> producerKafkaMerchant
   producerAPI - blueline >> producerKafkaCampaign
   producerAPI - blueline >> producerKafkaSpent
   producerAPI - bidirectional_edge >> producerDb
   producerAPI - bidirectional_edge >> producerLocalCache