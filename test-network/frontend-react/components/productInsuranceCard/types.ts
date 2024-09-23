type ButtonProps = {
	onClick: () => void;
	title: string;
};

type BadgeProps = {
	backgroundColor: string;
	textColor: string;
	text: string;
};

type ClaimStatus =
	| "Active"
	| "Pending"
	| "EvidencesApproved"
	| "EvidencesRejected"
	| "Approved"
	| "Rejected";

export type { BadgeProps, ClaimStatus, ButtonProps };
