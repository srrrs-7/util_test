<?php

namespace App\Jobs;

use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;

class BatchJob implements ShouldQueue
{
    use Queueable;

    /**
     * Create a new job instance.
     */
    public function __construct()
    {
        //
    }

    /**
     * Execute the job.
     */
    public function handle(): void
    {
        // Log the job completion
        \Log::info('Batch job completed successfully.');
        \Log::info('Batch job completed successfully.');
        \Log::info('Batch job completed successfully.');
        \Log::info('Batch job completed successfully.');
        // $this->delete(); // Delete the job from the queue
        // $this->delete(); // Delete the job from the queue
        // $this->delete(); // Delete the job from the queue
    }
}
