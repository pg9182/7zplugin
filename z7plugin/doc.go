// Package z7plugin implements 7-Zip's plugin interface.
//
// Must be built with `-buildmode=c-shared`. and the architecture must match
// 7-Zip. Put it in the `Formats“ directory in the install location. If the
// library doesn't load for some reason, try clicking on About in 7zFM, which
// will probably display an odd error message of some kind.
package z7plugin
