import { PageHeader } from "@/shared/ui/layout/page-header";

type ModulePlaceholderPageProps = {
  title: string;
  description?: string;
};

/**
 * Lightweight stand-in for modules that are not implemented yet.
 */
export function ModulePlaceholderPage({
  title,
  description = "This section is ready for your feature implementation. Connect routes to real screens under src/features.",
}: ModulePlaceholderPageProps) {
  return (
    <div className="space-y-6">
      <PageHeader title={title} description={description} />
    </div>
  );
}
