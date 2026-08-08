package demovisitor

// Package demovisitor documents the Demo / unknown visitor cookie pattern.
// Full middleware lives in product repos until extracted here.
//
// Cookie names: mm_demo, im_demo, etc.
// Hard rule: demo never mutates live deploy / stress / customer paths.

const Notice = "Demo / fake mode — not a live system"
