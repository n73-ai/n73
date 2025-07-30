export default function CoolText() {
  return (
    <div className="font-sans grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20">
      <h1 className="scroll-m-20 text-center text-4xl font-extrabold tracking-tight text-balance">
        Taxing Laughter: The Joke Tax Chronicles
      </h1>

      <h2 className="scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0">
        The People of the Kingdom
      </h2>

      <h3 className="scroll-m-20 text-2xl font-semibold tracking-tight">
        The Joke Tax
      </h3>

      <h4 className="scroll-m-20 text-xl font-semibold tracking-tight">
        People stopped telling jokes
      </h4>

      <p className="leading-7 [&:not(:first-child)]:mt-6">
        The king, seeing how much happier his subjects were, realized the error
        of his ways and repealed the joke tax.
      </p>

      <blockquote className="mt-6 border-l-2 pl-6 italic">
        &quot;After all,&quot; he said, &quot;everyone enjoys a good joke, so
        it&apos;s only fair that they should pay for the privilege.&quot;
      </blockquote>

      <ul className="my-6 ml-6 list-disc [&>li]:mt-2">
        <li>1st level of puns: 5 gold coins</li>
        <li>2nd level of jokes: 10 gold coins</li>
        <li>3rd level of one-liners : 20 gold coins</li>
      </ul>

      <code className="bg-muted relative rounded px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold">
        @radix-ui/react-alert-dialog
      </code>

      <p className="text-muted-foreground text-xl">
        A modal dialog that interrupts the user with important content and
        expects a response.
      </p>

      <div className="text-lg font-semibold">Are you absolutely sure?</div>

      <small className="text-sm leading-none font-medium">Email address</small>

      <p className="text-muted-foreground text-sm">Enter your email address.</p>
    </div>
  );
}
