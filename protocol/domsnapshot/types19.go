// +build go1.9

// Code generated by cdpgen. DO NOT EDIT.

package domsnapshot

import (
	"github.com/mafredri/cdp/protocol/dom"
	"github.com/mafredri/cdp/protocol/domdebugger"
	"github.com/mafredri/cdp/protocol/page"
)

// DOMNode A Node in the DOM tree.
type DOMNode struct {
	NodeType              int                         `json:"nodeType"`                        // `Node`'s nodeType.
	NodeName              string                      `json:"nodeName"`                        // `Node`'s nodeName.
	NodeValue             string                      `json:"nodeValue"`                       // `Node`'s nodeValue.
	TextValue             *string                     `json:"textValue,omitempty"`             // Only set for textarea elements, contains the text value.
	InputValue            *string                     `json:"inputValue,omitempty"`            // Only set for input elements, contains the input's associated text value.
	InputChecked          *bool                       `json:"inputChecked,omitempty"`          // Only set for radio and checkbox input elements, indicates if the element has been checked
	OptionSelected        *bool                       `json:"optionSelected,omitempty"`        // Only set for option elements, indicates if the element has been selected
	BackendNodeID         dom.BackendNodeID           `json:"backendNodeId"`                   // `Node`'s id, corresponds to DOM.Node.backendNodeId.
	ChildNodeIndexes      []int                       `json:"childNodeIndexes,omitempty"`      // The indexes of the node's child nodes in the `domNodes` array returned by `getSnapshot`, if any.
	Attributes            []NameValue                 `json:"attributes,omitempty"`            // Attributes of an `Element` node.
	PseudoElementIndexes  []int                       `json:"pseudoElementIndexes,omitempty"`  // Indexes of pseudo elements associated with this node in the `domNodes` array returned by `getSnapshot`, if any.
	LayoutNodeIndex       *int                        `json:"layoutNodeIndex,omitempty"`       // The index of the node's related layout tree node in the `layoutTreeNodes` array returned by `getSnapshot`, if any.
	DocumentURL           *string                     `json:"documentURL,omitempty"`           // Document URL that `Document` or `FrameOwner` node points to.
	BaseURL               *string                     `json:"baseURL,omitempty"`               // Base URL that `Document` or `FrameOwner` node uses for URL completion.
	ContentLanguage       *string                     `json:"contentLanguage,omitempty"`       // Only set for documents, contains the document's content language.
	DocumentEncoding      *string                     `json:"documentEncoding,omitempty"`      // Only set for documents, contains the document's character set encoding.
	PublicID              *string                     `json:"publicId,omitempty"`              // `DocumentType` node's publicId.
	SystemID              *string                     `json:"systemId,omitempty"`              // `DocumentType` node's systemId.
	FrameID               *page.FrameID               `json:"frameId,omitempty"`               // Frame ID for frame owner elements and also for the document node.
	ContentDocumentIndex  *int                        `json:"contentDocumentIndex,omitempty"`  // The index of a frame owner element's content document in the `domNodes` array returned by `getSnapshot`, if any.
	ImportedDocumentIndex *int                        `json:"importedDocumentIndex,omitempty"` // Index of the imported document's node of a link element in the `domNodes` array returned by `getSnapshot`, if any.
	TemplateContentIndex  *int                        `json:"templateContentIndex,omitempty"`  // Index of the content node of a template element in the `domNodes` array returned by `getSnapshot`.
	PseudoType            *dom.PseudoType             `json:"pseudoType,omitempty"`            // Type of a pseudo element node.
	ShadowRootType        *dom.ShadowRootType         `json:"shadowRootType,omitempty"`        // Shadow root type.
	IsClickable           *bool                       `json:"isClickable,omitempty"`           // Whether this DOM node responds to mouse clicks. This includes nodes that have had click event listeners attached via JavaScript as well as anchor tags that naturally navigate when clicked.
	EventListeners        []domdebugger.EventListener `json:"eventListeners,omitempty"`        // Details of the node's event listeners, if any.
	CurrentSourceURL      *string                     `json:"currentSourceURL,omitempty"`      // The selected url for nodes with a srcset attribute.
}
