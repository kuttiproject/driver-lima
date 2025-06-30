// Package driverlima implements a kutti driver for lima.
// It uses the limactl CLI to talk to lima.
//
// For cluster networking, it uses the lima user-v2 network.
//
// For nodes, it creates virtual machines from pre-built
// "cloud" images, maintained by the companion
// driver-lima-images project. These images are directly
// passed to lima for caching/VM disk creation.
//
// Thus, the driver itself only downloads the list of
// available images from the URL pointed to by the
// ImagesSourceURL variable.
package driverlima
