module: "namespacelabs.dev/foundation"
requirements: {
	api:          44
	toolsVersion: 4
}
prebuilts: {
	digest: {
		"namespacelabs.dev/foundation/internal/sdk/buf/image/prebuilt":                                "sha256:5a9f9711fcd93aa2cdb5d2ee2aaa1b2fdd23b7e139ff4a39438153668b9b84ef"
		"namespacelabs.dev/foundation/std/development/filesync/controller":                            "sha256:41ffa681aec6a70dcd5a7ebeccd94814688389a45f39810138a4d3f1ef8278da"
		"namespacelabs.dev/foundation/std/grpc/httptranscoding/configure":                             "sha256:eae106958cb4886d5ffab9f661f758a242d31c8c02358dd96fc301a247a7328a"
		"namespacelabs.dev/foundation/std/monitoring/grafana/tool":                                    "sha256:346a38e8301ba8366659280249a16ec287a14559f2855f5e7f2d07e5e4c190f9"
		"namespacelabs.dev/foundation/std/monitoring/prometheus/tool":                                 "sha256:067f86f8231c4787fa49d70251dba1c3b25d98bcfa020d21529994896786b5eb"
		"namespacelabs.dev/foundation/std/networking/gateway/controller":                              "sha256:11ff24b7079bd83001568570ccfac7b6118baa84f585901d54419bb7f08727a5"
		"namespacelabs.dev/foundation/std/networking/gateway/server/configure":                        "sha256:a6b6fcb1f42e730004aa0fdf339130dea9665df1a2581f517b78137bbb3631c7"
		"namespacelabs.dev/foundation/std/runtime/kubernetes/kube-state-metrics/configure":            "sha256:159e5af8e9c2724a272f1ff22a4d1b8d9e4f93e75fc8ac9b85309e36b6c8f676"
		"namespacelabs.dev/foundation/std/secrets/kubernetes":                                         "sha256:7c5be17536bdb60b99ef658dbc77f4f9ebbcbf10f343a79ce59aa2a9d91c9845"
		"namespacelabs.dev/foundation/std/startup/testdriver":                                         "sha256:39531c5b96518cee0a26037cb1ec7984a849d2f0a144ebf58c990832bdb5c9b0"
		"namespacelabs.dev/foundation/std/web/http/configure":                                         "sha256:128c028ef235bc9a2a2cd3ecce42298a4414b29acbddf1755f1f1c0014a927f5"
		"namespacelabs.dev/foundation/universe/aws/irsa/prepare":                                      "sha256:22f60c1f15911439a4711945245317acfa246184f94e2b7b5956131008c5dfe8"
		"namespacelabs.dev/foundation/universe/aws/s3/internal/configure":                             "sha256:0f2760d58ee3d4ec8aee1bd47d24d25cd730b888af047bf11ef21db570fff01d"
		"namespacelabs.dev/foundation/universe/aws/s3/internal/managebuckets/init":                    "sha256:40669c96749271e2f1247d98836d335949145415bda706c19bf6095a4a6df5f2"
		"namespacelabs.dev/foundation/universe/db/maria/incluster/tool":                               "sha256:252f83abd974d39c6ba258d21927dec1b514f893824f44ec7d7f0dc6e54e6b92"
		"namespacelabs.dev/foundation/universe/db/maria/internal/init":                                "sha256:1206bced820ab30286a5b3ad9baacbe1447e86e7aed4d2f2d2278fc0fa8a235a"
		"namespacelabs.dev/foundation/universe/db/maria/server/creds/tool":                            "sha256:0b0556ccca9e7e31d4e71779d5b9f4db7110f0b0f66593d0b0273f44b56e185e"
		"namespacelabs.dev/foundation/universe/db/maria/server/img":                                   "sha256:a7d5d37fe08eca6e91f88232784c92a6d411331a53aac7fcccb3b322875f9cb4"
		"namespacelabs.dev/foundation/universe/db/postgres/incluster/tool":                            "sha256:6495b69b8bea23dd7f7121229e15d0f785b2f9a92ec26fc2736c669cb474f089"
		"namespacelabs.dev/foundation/universe/db/postgres/internal/init":                             "sha256:d4b1e34b623bcb4e0ff535bab3725ff951f026b62b9f6a2016a78f279a195fed"
		"namespacelabs.dev/foundation/universe/db/postgres/opaque/tool":                               "sha256:42d380b51576ea211b7e9744dae5f0a960e20879a2da170cf7af538f15a24ede"
		"namespacelabs.dev/foundation/universe/db/postgres/rds/init":                                  "sha256:a3d4e0f1632d686bd0df6f573d1606a98ad324d2d2d1a217d155a117157bcbb2"
		"namespacelabs.dev/foundation/universe/db/postgres/rds/prepare":                               "sha256:52d9dccc910ded20216594218797a85b9810a4ca625087ac6071520b34d7a544"
		"namespacelabs.dev/foundation/universe/db/postgres/server/creds/tool":                         "sha256:205d7b5200e5f606ac5654951a5f2db0e1b6ca4ec55dc795cb9d780c29474ba9"
		"namespacelabs.dev/foundation/universe/db/postgres/server/img":                                "sha256:73b5e1ad8011ac702d6232ac5d3e10cb05a5baccc6ea3e87bbd79008f140eca0"
		"namespacelabs.dev/foundation/universe/development/localstack/s3/internal/configure":          "sha256:b72a8f03cb49e98c0d7c105086502e193a039b61dd5b7d30b8a06fc5bec9e71f"
		"namespacelabs.dev/foundation/universe/development/localstack/s3/internal/managebuckets/init": "sha256:59a43cac29183cb5df7bd6e61e2fe9ea6a3a582181f690fad8c5323bd3408037"
		"namespacelabs.dev/foundation/universe/networking/k8s-event-exporter/configure":               "sha256:44409819476881e3ed2e962fb3a3214500250495fc41ddb286a9503613dc091a"
		"namespacelabs.dev/foundation/universe/networking/tailscale/image":                            "sha256:444639fe064c0be98ddf66671d93db47ba973ab17636254906b228d69d5b06a4"
		"namespacelabs.dev/foundation/universe/storage/s3/internal/managebuckets":                     "sha256:595779e09f0b3f614b9b022489f6a6d4b6c6ceec894e5273cfe69bb9aadbe347"
		"namespacelabs.dev/foundation/universe/storage/s3/internal/prepare":                           "sha256:eead2b98bd0ff36c110c94441c61d7112c03432baf2ff5d51ddae1caa4d93db2"
	}
	baseRepository: "us-docker.pkg.dev/foundation-344819/prebuilts/"
}
internalAliases: [{
	module_name: "library.namespace.so"
	rel_path:    "library"
}]
