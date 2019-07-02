states = ["unmonitored", "stopped", "starting", "stopping", "restarting"]

initial_state = "unmonitored"

event "monitor" {
  transitions = [
    {
      from = ["unmonitored"]
      to   = "stopped"
    }
  ]
}

event "restart" {
  transitions = [
    {
      from = ["running", "stopped"]
      to   = "restarting"
    }
  ]
}

event "start" {
  transitions = [
    {
      from = ["unmonitored", "stopped"]
      to   = "starting"
    }
  ]
}

event "stop" {
  transitions = [
    {
      from = ["running"]
      to   = "stopping"
    }
  ]
}

event "tick" {
  transitions = [
    {
      from     = ["starting"]
      to       = "running"
      if_guard = {
        registered_func = "is-process-running"
      }
    },
    {
      from         = ["starting"]
      to           = "stopped"
      unless_guard = {
        registered_func = "is-process-running"
      }
    },
    {
      from         = ["running"]
      to           = "stopped"
      unless_guard = {
        registered_func = "is-process-running"
      }
    },
    {
      from     = ["stopping"]
      to       = "running"
      if_guard = {
        registered_func = "is-process-running"
      }
    },
    {
      from         = ["stopping"]
      to           = "stopped"
      unless_guard = {
        registered_func = "is-process-running"
      }
    },
    {
      from     = ["stopped"]
      to       = "running"
      if_guard = {
        registered_func = "is-process-running"
      }
    },
    {
      from         = ["stopped"]
      to           = "starting"
      if_guard     = {
        registered_func = "is-autostart-on"
      }
      unless_guard = {
        registered_func = "is-process-running"
      }
    },
    {
      from     = ["restarting"]
      to       = "running"
      if_guard = {
        registered_func = "is-process-running"
      }
    },
    {
      from         = ["restarting"]
      to           = "stopped"
      unless_guard = {
        registered_func = "is-process-running"
      }
    }
  ]
}

event "unmonitor" {
  transitions = [
    {
      from = []
      to   = "unmonitored"
    }
  ]
}

submachine running {
  states = ["pending", "success", "failure"]

  initial_state = "pending"

  event "process" {
    transitions = [
      {
        from = ["pending"]
        to   = "processing"
      }
    ]
  }

  event "succeed" {
    transitions = [
      {
        from = ["processing"]
        to   = "success"
      }
    ]
  }

  event "fail" {
    transitions = [
      {
        from = ["processing"]
        to   = "failure"
      }
    ]
  }

  submachine processing {
    states = ["loading", "subsubprocessing"]

    initial_state = "loading"

    event "subsubprocess" {
      transitions = [
        {
          from = ["loading"]
          to   = "subsubprocessing"
        }
      ]
    }

    event "to_done" {
      transitions = [
        {
          from = ["subsubprocessing"]
          to   = "done"
        }
      ]
    }

    around_callbacks = [
      {
        do = {
          registered_func = "subsub-around-callback-1"
        }
      }
    ]

    after_callbacks = [
      {
        to = ["subsubprocessing"]
        do = {
          registered_func = "subsub-after-callback-1"
        }
      },
      {
        to        = ["done"]
        exit_into = "success"
      }
    ]

    failure_callbacks = [
      {
        do = {
          registered_func = "subsub-failure-callback-1"
        }
      }
    ]
  }
  around_callbacks = [
    {
      do = {
        registered_func = "sub-around-callback-1"
      }
    }
  ]

  after_callbacks = [
    {
      to = ["processing"]

      do = {
        registered_func = "sub-after-callback-1"
      }
    },
    {
      to        = ["success"]
      exit_into = "stopped"
    },
    {
      to        = ["failure"]
      exit_into = "retrying"
    }
  ]

  failure_callbacks = [
    {
      do = {
        registered_func = "sub-failure-callback-1"
      }
    }
  ]
}

before_callbacks = [
  {
    to = ["starting"]
    do = {
      registered_func = "before-callback-1"
    }
  },
  {
    to = ["stopping"]
    do = {
      registered_func = "before-callback-2"
    }
  },
  {
    to = ["restarting"]
    do = {
      registered_func = "before-callback-3"
    }
  },
  {
    to = ["unmonitored"]
    do = {
      registered_func = "before-callback-4"
    }
  }
]

around_callbacks = [
  {
    do = {
      registered_func = "around-callback-"
    }
  }
]

after_callbacks = [
  {
    to = ["starting"]
    do = {
      registered_func = "after-callback-1"
    }
  },
  {
    to = ["stopping"]
    do = {
      registered_func = "after-callback-2"
    }
  },
  {
    to = ["restarting"]
    do = {
      registered_func = "after-callback-3"
    }
  },
  {
    to = ["running"]
    do = {
      registered_func = "after-callback-4"
    }
  }
]

failure_callbacks = [
  {
    do = {
      registered_func = "failure-callback"
    }
  }
]
