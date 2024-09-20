type BadgeProps = {
	backgroundColor: string;
	textColor: string;
	text: string;
};

type ClaimStatus = "Active" | "Pending" | "Approved" | "Rejected";

export type { BadgeProps, ClaimStatus };
